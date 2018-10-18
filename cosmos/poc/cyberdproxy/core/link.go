package core

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cybercongress/cyberd/cosmos/poc/app"
	"io/ioutil"
	"net/http"
)

type LinkRequest struct {
	Fee        auth.StdFee   `json:"fee"`
	Msgs       []app.MsgLink `json:"msgs"`
	Signatures []Signature   `json:"signatures"`
	Memo       string        `json:"memo"`
}

func LinkHandlerFn(ctx ProxyContext) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		requestBytes, _ := ioutil.ReadAll(r.Body)
		var request LinkRequest
		json.Unmarshal(requestBytes, &request)

		// BUILDING COSMOS SDK TX
		signatures := make([]auth.StdSignature, 0, len(request.Signatures))
		for _, sig := range request.Signatures {
			stdSig := auth.StdSignature{
				PubKey: sig.PubKey, Signature: sig.Signature, AccountNumber: sig.AccountNumber, Sequence: sig.Sequence,
			}
			signatures = append(signatures, stdSig)
		}

		msgs := make([]sdk.Msg, 0, len(request.Msgs))
		for _, msg := range request.Msgs {
			msgs = append(msgs, msg)
		}

		stdTx := auth.StdTx{Msgs: msgs, Fee: request.Fee, Signatures: signatures, Memo: request.Memo}

		stdTxBytes, _ := ctx.Codec.MarshalBinary(stdTx)

		resp, _ := ctx.Node.BroadcastTxCommit(stdTxBytes)

		respBytes, _ := json.Marshal(resp)
		w.Write(respBytes)
	}
}