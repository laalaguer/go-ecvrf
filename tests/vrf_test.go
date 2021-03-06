// Copyright (c) 2020 vechain.org.
// Licensed under the MIT license.

package tests

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"

	"math/big"
	"reflect"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/vechain/go-ecvrf"
)

// Case Testing cases structure.
type Case struct {
	Sk    string `json:"sk"`
	Pk    string `json:"pk"`
	Alpha string `json:"alpha"`
	Pi    string `json:"pi"`
	Beta  string `json:"beta"`
}

func readCases(fileName string) ([]Case, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err2 := ioutil.ReadAll(jsonFile)
	if err2 != nil {
		return nil, err2
	}

	var cases = make([]Case, 0)
	err3 := json.Unmarshal(byteValue, &cases)
	if err3 != nil {
		return cases, err3
	}

	return cases, nil
}

func Test_Secp256K1Sha256Tai_vrf_Prove(t *testing.T) {
	// Know Correct cases.
	var cases, _ = readCases("./secp256_k1_sha256_tai.json")

	type Test struct {
		name     string
		sk       *ecdsa.PrivateKey
		alpha    []byte
		wantBeta []byte
		wantPi   []byte
		wantErr  bool
	}

	tests := []Test{}
	for _, c := range cases {
		skBytes, _ := hex.DecodeString(c.Sk)
		sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)

		alpha, _ := hex.DecodeString(c.Alpha)
		wantBeta, _ := hex.DecodeString(c.Beta)
		wantPi, _ := hex.DecodeString(c.Pi)

		tests = append(tests, Test{
			c.Sk,
			sk.ToECDSA(),
			alpha,
			wantBeta,
			wantPi,
			false,
		})
	}

	vrf := ecvrf.NewSecp256k1Sha256Tai()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vrf
			gotBeta, gotPi, err := v.Prove(tt.sk, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf.Prove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBeta, tt.wantBeta) {
				t.Errorf("vrf.Prove() gotBeta = %v, want %v", hex.EncodeToString(gotBeta), hex.EncodeToString(tt.wantBeta))
			}
			if !reflect.DeepEqual(gotPi, tt.wantPi) {
				t.Errorf("vrf.Prove() gotPi = %v, want %v", hex.EncodeToString(gotPi), hex.EncodeToString(tt.wantPi))
			}
		})
	}
}

func Test_Secp256K1Sha256Tai_vrf_Verify(t *testing.T) {
	// Know Correct cases.
	var cases, _ = readCases("./secp256_k1_sha256_tai.json")

	type Test struct {
		name     string
		pk       *ecdsa.PublicKey
		alpha    []byte
		pi       []byte
		wantBeta []byte
		wantErr  bool
	}

	tests := []Test{}
	for _, c := range cases {
		skBytes, _ := hex.DecodeString(c.Sk)
		sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), skBytes)

		pk := sk.PubKey().ToECDSA()

		alpha, _ := hex.DecodeString(c.Alpha)

		wantPi, _ := hex.DecodeString(c.Pi)

		wantBeta, _ := hex.DecodeString(c.Beta)

		tests = append(tests, Test{
			c.Alpha,
			pk,
			alpha,
			wantPi,
			wantBeta,
			false,
		})
	}

	vrf := ecvrf.NewSecp256k1Sha256Tai()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vrf
			gotBeta, err := v.Verify(tt.pk, tt.alpha, tt.pi)
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBeta, tt.wantBeta) {
				t.Errorf("vrf.Verify() = %v, want %v", gotBeta, tt.wantBeta)
			}
		})
	}
}

func Test_P256Sha256Tai_vrf_Prove(t *testing.T) {
	// Know Correct cases.
	var P256Sha256TaiCases, _ = readCases("./p256_sha256_tai.json")

	type Test struct {
		name     string
		sk       *ecdsa.PrivateKey
		alpha    []byte
		wantBeta []byte
		wantPi   []byte
		wantErr  bool
	}

	tests := []Test{}
	for _, c := range P256Sha256TaiCases {
		skBytes, _ := hex.DecodeString(c.Sk)
		curve := elliptic.P256()
		pkX, pkY := curve.ScalarBaseMult(skBytes)
		sk := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: curve,
				X:     pkX,
				Y:     pkY,
			},
			D: new(big.Int).SetBytes(skBytes),
		}
		alpha, _ := hex.DecodeString(c.Alpha)
		wantBeta, _ := hex.DecodeString(c.Beta)
		wantPi, _ := hex.DecodeString(c.Pi)

		tests = append(tests, Test{
			c.Alpha,
			sk,
			alpha,
			wantBeta,
			wantPi,
			false,
		})
	}

	vrf := ecvrf.NewP256Sha256Tai()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vrf
			gotBeta, gotPi, err := v.Prove(tt.sk, tt.alpha)
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf.Prove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBeta, tt.wantBeta) {
				t.Errorf("vrf.Prove() gotBeta = %v, want %v", gotBeta, tt.wantBeta)
			}
			if !reflect.DeepEqual(gotPi, tt.wantPi) {
				t.Errorf("vrf.Prove() gotPi = %v, want %v", gotPi, tt.wantPi)
			}
		})
	}
}

func Test_P256Sha256Tai_vrf_Verify(t *testing.T) {
	// Know Correct cases.
	var P256Sha256TaiCases, _ = readCases("./p256_sha256_tai.json")

	type Test struct {
		name     string
		pk       *ecdsa.PublicKey
		alpha    []byte
		pi       []byte
		wantBeta []byte
		wantErr  bool
	}

	tests := []Test{}
	for _, c := range P256Sha256TaiCases {
		curve := elliptic.P256()
		skBytes, _ := hex.DecodeString(c.Sk)

		pkX, pkY := curve.ScalarBaseMult(skBytes)
		pk := ecdsa.PublicKey{
			Curve: curve,
			X:     pkX,
			Y:     pkY,
		}

		alpha, _ := hex.DecodeString(c.Alpha)
		pi, _ := hex.DecodeString(c.Pi)
		wantBeta, _ := hex.DecodeString(c.Beta)

		tests = append(tests, Test{
			c.Alpha,
			&pk,
			alpha,
			pi,
			wantBeta,
			false,
		})
	}

	vrf := ecvrf.NewP256Sha256Tai()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vrf
			gotBeta, err := v.Verify(tt.pk, tt.alpha, tt.pi)
			if (err != nil) != tt.wantErr {
				t.Errorf("vrf.Verify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBeta, tt.wantBeta) {
				t.Errorf("vrf.Verify() = %v, want %v", gotBeta, tt.wantBeta)
			}
		})
	}
}
