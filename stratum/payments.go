// Many payments integration functions and ideas from: https://github.com/JKKGBE/open-zcash-pool which is a fork of https://github.com/sammy007/open-ethereum-pool
package stratum

import (
	"fmt"
	"log"
	"math/big"
	"time"

	"git.dero.io/Nelbert442/dero-golang-pool/pool"
	"git.dero.io/Nelbert442/dero-golang-pool/rpc"
	"git.dero.io/Nelbert442/dero-golang-pool/util"
	"github.com/deroproject/derosuite/globals"
	"github.com/deroproject/derosuite/walletapi"
)

type PayoutsProcessor struct {
	config   *pool.PaymentsConfig
	rpc      *rpc.RPCClient
	halt     bool
	lastFail error
}

type PayoutTracker struct {
	Destinations []rpc.Destinations
	PaymentIDs   []string
}

/* Used when integrating with derosuite functions, currently not being used but in place incase functions are used later
type Transfer struct {
	rAddress	*address.Address
	PaymentID	[]byte
	Amount		uint64
	Fees		uint64
	TX			*transaction.Transaction
	TXID		crypto.Hash
	Size		float32
	Status		string
	Inputs		[]uint64
	InputSum	uint64
	Change		uint64
	Relay		bool
	OfflineTX	bool
	Filename	string
}
*/

var wallet *walletapi.Wallet

func NewPayoutsProcessor(cfg *pool.PaymentsConfig, s *StratumServer) *PayoutsProcessor {
	u := &PayoutsProcessor{config: cfg} //backend: s.backend}
	u.rpc = s.rpc()
	return u
}

func (u *PayoutsProcessor) Start(s *StratumServer) {
	log.Println("[Payments] Starting payouts")

	intv, _ := time.ParseDuration(u.config.Interval)
	timer := time.NewTimer(intv)
	log.Printf("[Payments] Set payouts interval to %v", intv)

	payments := Graviton_backend.GetPendingPayments()

	if len(payments) > 0 {
		log.Printf("[Payments] Previous payout failed, trying to resolve it. List of failed payments:\n %v",
			formatPendingPayments(payments))
	}

	// Immediately process payouts after start
	u.process(s)
	timer.Reset(intv)

	go func() {
		for {
			select {
			case <-timer.C:
				u.process(s)
				timer.Reset(intv)
			}
		}
	}()
}

func (u *PayoutsProcessor) process(s *StratumServer) {

	maxAddresses := u.config.MaxAddresses
	var payoutList []rpc.Destinations
	var paymentIDPayeeList []rpc.Destinations
	var payIDList []string
	var payPending []*PaymentPending
	//paymentsToRemove := make(map[string][]int)
	//paymentsToRemove := make(map[string]int)

	walletURL := fmt.Sprintf("http://%s:%v/json_rpc", u.config.WalletHost, u.config.WalletPort)
	mustPay := 0
	minersPaid := 0
	totalAmount := big.NewInt(0)

	// Graviton DB Pending Balance
	payPending = Graviton_backend.GetPendingPayments()
	for _, val := range payPending {

		login := val.Address
		amount := val.Amount

		if !u.reachedThreshold(amount) {
			continue
		}

		// Add index of payPending to array of paymentsToRemove which will be used later to cleanup DB
		//paymentsToRemove[login] = i

		// Check if we have enough funds
		poolBalanceObj, err := u.rpc.GetBalance(walletURL)
		if err != nil {
			// TODO: mark sick maybe for tracking and frontend reporting?
			log.Printf("[Payments] Error when getting balance from wallet %s. Will try again in %s", walletURL, u.config.Interval)
			break
		}
		poolBalance := poolBalanceObj.UnlockedBalance
		if poolBalance < amount {
			log.Printf("[Payments] Not enough balance for payment, need %v DERO, pool has %v DERO", amount, poolBalance)
			break
		}

		// Address validations
		// We already do this for when the miner connects, we need to get those details/vars or just regen them as well as re-validate JUST TO BE SURE prior to attempting to send
		// NOTE: The issue with grabbing from the miners arr (s.miners), is that if they're not actively mining but get rewards from a past round, the query will not return their detail for payout

		address, _, paymentID, _, _ := s.splitLoginString(login)

		log.Printf("[Payments] Split login. Address: %v, paymentID: %v", address, paymentID)

		// Validate Address
		validatedAddress, err := globals.ParseValidateAddress(address)
		_ = validatedAddress

		if err != nil {
			log.Printf("[Payments] Invalid address format. Will not process payments - %v", address)
			continue
		}

		currAddr := rpc.Destinations{
			Amount:  amount,
			Address: address,
		}
		mustPay++

		// If paymentID, put in an array that'll be walked through one at a time versus combining addresses/amounts.
		if paymentID != "" {
			paymentIDPayeeList = append(paymentIDPayeeList, currAddr)
			payIDList = append(payIDList, paymentID)
		} else {
			payoutList = append(payoutList, currAddr)
		}
	}

	payIDTracker := &PayoutTracker{Destinations: paymentIDPayeeList, PaymentIDs: payIDList}

	// If there are payouts to be processed
	if len(payoutList) > 0 || len(payIDList) > 0 {

		// Send DERO - Native, TODO)
		// Issues with mutex and also running wallet process locally (with --rpc-server), as it can't get a lock to generate the transaction.
		// Problem comes in, then it can't use the daemon/wallet to do the work it needs and errors, since we're not providing authentication in config etc.
		// Future may be allow for auth to rpc-server within the config.json or other means as an option to run it that way. Otherwise will continue using the rpc option as the future option.
		/*
			// Validate Address
			wallet.SetDaemonAddress("http://127.0.0.1:30306")

			transfer.rAddress, err = globals.ParseValidateAddress(login)
			if err != nil {
				log.Printf("Invalid address format - %v", login)
				break
			} else {
				log.Printf("Valid Address")
			}

			transfer.PaymentID = nil

			// This fails out w/ mutex error, not sure TODO
			//curBalance, _ := wallet.Get_Balance()

			transfer.Amount = amount
			addr_list := []address.Address{*transfer.rAddress}
			amount_list := []uint64{transfer.Amount}
			fees_per_kb := uint64(0) // fees  must be calculated by walletapi

			tx, inputs, input_sum, change, err := wallet.Transfer(addr_list, amount_list, 0, hex.EncodeToString(transfer.PaymentID), fees_per_kb, 0)
			_ = inputs

			if err != nil {
				log.Printf("Error while building transaction: %s", err)
				break
			}

			transfer.OfflineTX = false

			transfer.Relay = build_relay_transaction(tx, inputs, input_sum, change, err, transfer.OfflineTX, amount_list)

			if !transfer.Relay {
				log.Printf("Error:  Unable to build the transfer.")
				break
			}

			err = wallet.SendTransaction(transfer.TX) // relay tx to daemon/network

			if err == nil {
				transfer.Status = "Success"
				transfer.TXID = transfer.TX.GetHash()
			} else {
				transfer.Status = "Failed"
				transfer.TXID = transfer.TX.GetHash()
				log.Printf("Error relaying transaction: %s", err)
				break
			}
			txHash := transfer.TXID.String()
		*/

		// Send DERO - RPC (working)
		var currPayout rpc.Transfer_Params
		var lastPos int
		currPayout.Mixin = u.config.Mixin
		currPayout.Unlock_time = 0
		currPayout.Get_tx_key = true
		currPayout.Do_not_relay = false
		currPayout.Get_tx_hex = true

		// Payout paymentID addresses, one at a time since paymentID is used in the tx generation and is a non-array input
		for p, payee := range payIDTracker.Destinations {
			currPayout.Payment_ID = payIDTracker.PaymentIDs[p]
			currPayout.Destinations = append(currPayout.Destinations, payee)

			paymentOutput, err := u.rpc.SendTransaction(walletURL, currPayout)

			if err != nil {
				log.Printf("[Payments] Error with transaction: %v", err)
				break
			}
			log.Printf("[Payments] Success: %v", paymentOutput)
			// Log transaction hash
			txHash := paymentOutput.Tx_hash_list
			txFee := paymentOutput.Fee_list
			// As pool owner, you probably want to store keys so that you can prove a send if required.
			txKey := paymentOutput.Tx_key_list

			if txHash == nil {
				log.Printf("[Payments] Failed to generate transaction. It was sent successfully to rpc server, but no reply back.")

				break
			}

			// Debit miner's balance and update stats
			login := payee.Address + s.config.Stratum.PaymentID.AddressSeparator + currPayout.Payment_ID
			amount := payee.Amount

			// Remove pending payout from graviton db
			//log.Printf("[Payments] Before Payment Pruning: %v", payPending)
			//f := paymentsToRemove[login] - p
			//log.Printf("[Payments] Removing payPending: %v", payPending[f])
			//payPending = removePendingPayments(payPending, f)

			for j, f := range payPending {
				if login == f.Address {
					//log.Printf("[Payments] Removing payPending: %v", payPending[j])
					payPending = removePendingPayments(payPending, j)
					break
				}
			}
			/*for j, f := range paymentsToRemove[login] {
				// This keeps the array val in check while processing payments.
				// Example:
				// login 1: f = 0
				// login 2: f = 1
				// login 1 gets pruned, payPending length is now 1, need to reduce login 2 f = 0 instead of 1, else out of range index, so 1 - 1 and so on..
				if f <= j {
					f = j - f
				}
				payPending = removePendingPayments(payPending, f)
			}*/

			prunedPaymentsPending := &PendingPayments{PendingPayout: payPending}

			err = Graviton_backend.OverwritePendingPayments(prunedPaymentsPending)
			if err != nil {
				log.Printf("[Payments] Error overwriting pending payments. %v", err)
				break
			}

			// Update stats for pool payments (gravitondb)
			info := &MinerPayments{}
			info.Login = login
			info.TxHash = txHash[0]
			info.TxKey = txKey[0]
			info.TxFee = txFee[0]
			info.Mixin = u.config.Mixin
			info.Amount = amount
			info.Timestamp = util.MakeTimestamp() / 1000
			infoErr := s.gravitonDB.WriteProcessedPayments(info)
			if infoErr != nil {
				log.Printf("[Payments] Graviton DB err: %v", infoErr)
				break
			}

			minersPaid++
			totalAmount.Add(totalAmount, big.NewInt(int64(amount)))
			log.Printf("[Payments] Paid %v DERO to %v, PaymentID: %v, TxHash: %v, Fee: %v, Mixin: %v", int64(amount), login, currPayout.Payment_ID, txHash[0], txFee[0], u.config.Mixin)
		}
		currPayout.Destinations = nil

		// Payout non-paymentID addresses, max at a time according to maxAddresses in config
		for i, value := range payoutList {
			currPayout.Payment_ID = ""
			currPayout.Destinations = append(currPayout.Destinations, value)

			// Payout if maxAddresses is reached or the payout list ending is reached
			if len(currPayout.Destinations) >= int(maxAddresses) || i+1 == len(payoutList) {
				paymentOutput, err := u.rpc.SendTransaction(walletURL, currPayout)

				if err != nil {
					log.Printf("[Payments] Error with transaction: %v", err)
					break
				}
				log.Printf("[Payments] Success: %v", paymentOutput)
				// Log transaction hash
				txHash := paymentOutput.Tx_hash_list
				txFee := paymentOutput.Fee_list
				// As pool owner, you probably want to store keys so that you can prove a send if required.
				txKey := paymentOutput.Tx_key_list

				if txHash == nil {
					log.Printf("[Payments] Failed to generate transaction. It was sent successfully to rpc server, but no reply back.")

					break
				}

				if len(payoutList) > 1 && maxAddresses > 1 {
					payPos := i - lastPos
					for k := 0; k <= payPos; k++ {
						// Debit miner's balance and update stats
						login := payoutList[lastPos+k].Address
						amount := payoutList[lastPos+k].Amount

						// Remove pending payout from graviton db
						//log.Printf("[Payments] Before Payment Pruning: %v", payPending)
						//f := paymentsToRemove[login] - k
						//log.Printf("[Payments] Removing payPending: %v", payPending[f])
						//payPending = removePendingPayments(payPending, f)

						for j, f := range payPending {
							if login == f.Address {
								//log.Printf("[Payments] Removing payPending: %v", payPending[j])
								payPending = removePendingPayments(payPending, j)
								break
							}
						}
						/*
							for j, f := range paymentsToRemove[login] {
								// This keeps the array val in check while processing payments.
								// Example:
								// login 1: f = 0
								// login 2: f = 1
								// login 1 gets pruned, payPending length is now 1, need to reduce login 2 f = 0 instead of 1, else out of range index, so 1 - 1 and so on..
								log.Printf("f: %v; j: %v", f, j)
								if f <= j {
									f = j - f
								}
								log.Printf("f: %v; j: %v", f, j)
								payPending = removePendingPayments(payPending, f)
							}
						*/

						prunedPaymentsPending := &PendingPayments{PendingPayout: payPending}

						err = Graviton_backend.OverwritePendingPayments(prunedPaymentsPending)
						if err != nil {
							log.Printf("[Payments] Error overwriting pending payments. %v", err)
							break
						}

						// Update stats for pool payments (gravitondb)
						info := &MinerPayments{}
						info.Login = login
						info.TxHash = txHash[0]
						info.TxKey = txKey[0]
						info.TxFee = txFee[0]
						info.Mixin = u.config.Mixin
						info.Amount = amount
						info.Timestamp = util.MakeTimestamp() / 1000
						infoErr := s.gravitonDB.WriteProcessedPayments(info)
						if infoErr != nil {
							log.Printf("[Payments] Graviton DB err: %v", infoErr)
							break
						}

						minersPaid++
						totalAmount.Add(totalAmount, big.NewInt(int64(amount)))
						log.Printf("[Payments] Paid %v DERO to %v, TxHash: %v, Fee: %v, Mixin: %v", int64(amount), login, txHash[0], txFee[0], u.config.Mixin)
					}
				} else {
					log.Printf("[Payments] Processing payoutList[i]: %v", payoutList[i])
					// Debit miner's balance and update stats
					login := value.Address
					amount := value.Amount

					// Remove pending payout from graviton db
					//log.Printf("[Payments] Before Payment Pruning: %v", payPending)
					//f := paymentsToRemove[login]
					//log.Printf("[Payments] Removing payPending: %v", payPending[f])
					//payPending = removePendingPayments(payPending, f)

					for j, f := range payPending {
						if login == f.Address {
							//log.Printf("[Payments] Removing payPending: %v", payPending[j])
							payPending = removePendingPayments(payPending, j)
							break
						}
					}

					/*
						for j, f := range paymentsToRemove[login] {
							// This keeps the array val in check while processing payments.
							// Example:
							// login 1: f = 0
							// login 2: f = 1
							// login 1 gets pruned, payPending length is now 1, need to reduce login 2 f = 0 instead of 1, else out of range index, so 1 - 1 and so on..
							if f <= j {
								f = j - f
							}
							payPending = removePendingPayments(payPending, f)
						}
					*/

					prunedPaymentsPending := &PendingPayments{PendingPayout: payPending}

					err = Graviton_backend.OverwritePendingPayments(prunedPaymentsPending)
					if err != nil {
						log.Printf("[Payments] Error overwriting pending payments. %v", err)
						break
					}

					// Update stats for pool payments (gravitondb)
					info := &MinerPayments{}
					info.Login = login
					info.TxHash = txHash[0]
					info.TxKey = txKey[0]
					info.TxFee = txFee[0]
					info.Mixin = u.config.Mixin
					info.Amount = amount
					info.Timestamp = util.MakeTimestamp() / 1000
					infoErr := s.gravitonDB.WriteProcessedPayments(info)
					if infoErr != nil {
						log.Printf("[Payments] Graviton DB err: %v", infoErr)
						break
					}

					minersPaid++
					totalAmount.Add(totalAmount, big.NewInt(int64(amount)))
					log.Printf("[Payments] Paid %v DERO to %v, TxHash: %v, Fee: %v, Mixin: %v", int64(amount), login, txHash[0], txFee[0], u.config.Mixin)
				}
				// Empty currpayout destinations array
				currPayout.Destinations = nil
				lastPos = i + 1 // Increment lastPos so it'll be equal to i next round in loop (if required)
			}
		}
	}

	if mustPay > 0 {
		log.Printf("[Payments] Paid total %v DERO to %v of %v payees", totalAmount, minersPaid, mustPay)
	} else {
		log.Println("[Payments] No payees that have reached payout threshold")
	}
}

func removePendingPayments(s []*PaymentPending, i int) []*PaymentPending {
	if len(s) == 1 {
		return nil
	}
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func formatPendingPayments(list []*PaymentPending) string {
	var s string
	for _, v := range list {
		s += fmt.Sprintf("\t[Payments] Address: %s, Amount: %v DERO, %v\n", v.Address, v.Amount, time.Unix(v.Timestamp, 0))
	}
	return s
}

func (self PayoutsProcessor) reachedThreshold(amount uint64) bool {
	return self.config.Threshold < amount
}
