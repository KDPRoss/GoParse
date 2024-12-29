package main

// ========== Kinds ==========

// Transliteration of the OCaml type:
// type knd = KStar
//          | KArr  of (knd * knd)

type Kind interface{ kindEvidence() }
type KStar struct{}
type KArr struct{ K1, K2 Kind }

func (KStar) kindEvidence() {}
func (KArr) kindEvidence()  {}

// ========== Types ==========

// Transliteration of the OCaml type:
// type typ = TVar  of string
//          | TAbs  of (string * knd * typ)
//          | TCVal of (string * typ list)
//          | TArr  of (typ * typ)
//          | TApp  of (typ * typ)
//          | TTpl  of typ list

type Type interface{ typeEvidence() }
type TVar struct{ X string }
type TAbs struct {
	X string
	K Kind
	T Type
}
type TCVal struct {
	C  string
	Ts []Type
}
type TArr struct{ T1, T2 Type }
type TApp struct{ T1, T2 Type }
type TTpl struct{ Ts []Type }

func (TVar) typeEvidence()  {}
func (TAbs) typeEvidence()  {}
func (TCVal) typeEvidence() {}
func (TArr) typeEvidence()  {}
func (TApp) typeEvidence()  {}
func (TTpl) typeEvidence()  {}
