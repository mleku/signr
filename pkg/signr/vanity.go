package signr

import (
	"encoding/hex"
	"fmt"
	"mleku.online/git/atomic"
	"mleku.online/git/bech32"
	secp "mleku.online/git/ec/secp"
	"os"
	"strings"
	"sync"
	"time"

	"mleku.online/git/ec/schnorr"
	"mleku.online/git/interrupt"
	"mleku.online/git/qu"
	"mleku.online/git/signr/pkg/nostr"
)

type Position int

const prefix = nostr.PubHRP + "1"

const (
	PositionBeginning = iota
	PositionContains
	PositionEnding
)

type Result struct {
	sec  *secp.SecretKey
	npub string
	pub  *secp.PublicKey
}

func (s *Signr) Vanity(str, name string, where Position,
	threads int) (e error) {

	qu.SetLogging(true)
	// check the string has valid bech32 ciphers
	for i := range str {
		wrong := true
		for j := range bech32.Charset {
			if str[i] == bech32.Charset[j] {
				wrong = false
				break
			}
		}
		if wrong {
			return fmt.Errorf("found invalid character '%c' only ones from '%s' allowed\n",
				str[i], bech32.Charset)
		}
	}
	started := time.Now()
	quit, shutdown := qu.T(), qu.T()
	resC := make(chan Result)
	interrupt.AddHandler(func() {
		// this will stop work if CTRL-C or Interrupt signal from OS.
		shutdown.Q()
	})
	var wg sync.WaitGroup
	counter := atomic.NewInt64(0)
	for i := 0; i < threads; i++ {
		s.Log("starting up worker %d\n", i)
		go s.mine(str, where, quit, resC, &wg, counter)
	}
	tick := time.NewTicker(time.Second * 5)
	var res Result
out:
	for {
		select {
		case <-tick.C:
			workingFor := time.Now().Sub(started)
			wm := workingFor % time.Second
			workingFor -= wm
			s.Log("working for %v, attempt %d\n",
				workingFor, counter.Load())
		case r := <-resC:
			// one of the workers found the solution
			res = r
			// tell the others to stop
			quit.Q()
			break out
		case <-shutdown.Wait():
			quit.Q()
			s.Info("\ninterrupt signal received\n")
			os.Exit(0)
		}
	}

	// wait for all of the workers to stop
	wg.Wait()

	s.Info("generated in %d attempts, taking %v\n", counter.Load(),
		time.Now().Sub(started))
	secBytes := res.sec.Serialize()
	if s.Verbose.Load() {
		s.Log(
			"generated key pair:\n"+
				"\nhex:\n"+
				"\tsecret: %s\n"+
				"\tpublic: %s\n\n",
			hex.EncodeToString(secBytes),
			hex.EncodeToString(schnorr.SerializePubKey(res.pub)),
		)
		nsec, _ := nostr.SecretKeyToNsec(res.sec)
		s.Log("nostr:\n"+
			"\tsecret: %s\n"+
			"\tpublic: %s\n\n",
			nsec, res.npub)
	}
	if e = s.Save(name, secBytes, res.npub); e != nil {
		e = fmt.Errorf("error saving keys: %v", e)
		return
	}
	return
}

func (s *Signr) mine(str string, where Position,
	quit qu.C, resC chan Result, wg *sync.WaitGroup,
	counter *atomic.Int64) {

	wg.Add(1)
	var r Result
	var e error
	found := false
out:
	for {
		select {
		case <-quit:
			wg.Done()
			if found {
				// send back the result
				s.Log("sending back result\n")
				resC <- r
				s.Log("sent\n")
			} else {
				s.Log("other thread found it\n")
			}
			break out
		default:
		}
		counter.Inc()
		r.sec, r.pub, e = s.GenKeyPair()
		if e != nil {
			s.Err("error generating key: '%v' worker stopping", e)
			break out
		}
		r.npub, e = nostr.PublicKeyToNpub(r.pub)
		if e != nil {
			s.Err("fatal error generating npub: %s\n", e)
			break out
		}
		switch where {
		case PositionBeginning:
			if strings.HasPrefix(r.npub, prefix+str) {
				found = true
				s.Info("found it %s %d attempts\n", r.npub, counter.Load())
				quit.Q()
			}
		case PositionEnding:
			if strings.HasSuffix(r.npub, str) {
				found = true
				s.Info("found it %s %d attempts\n", r.npub, counter.Load())
				quit.Q()
			}
		case PositionContains:
			if strings.Contains(r.npub, str) {
				found = true
				s.Info("found it %s %d attempts\n", r.npub, counter.Load())
				quit.Q()
			}
		}
	}
}
