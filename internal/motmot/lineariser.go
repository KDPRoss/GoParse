package main

import "fmt"

func ShowKind(k Kind) string {
	switch kT := k.(type) {
	case KStar:
		return "*"
	case KArr:
		return fmt.Sprintf("(%s -> %s)", ShowKind(kT.K1), ShowKind(kT.K2))
	default:
		// Cannot happen because we've covered all `Kind` ADT
		// constructors.
		panic("impossible")
	}
}

func ShowType(t Type) string {
	switch tT := t.(type) {
	case TVar:
		return tT.X
	case TAbs:
		return fmt.Sprintf("((%s : %s) => %s)", tT.X, ShowKind(tT.K), ShowType(tT.T))
	case TCVal:
		{
			res := tT.C
			if len(tT.Ts) > 0 {
				for _, t := range tT.Ts {
					res = fmt.Sprintf("%s %s", res, ShowType(t))
				}

				res = fmt.Sprintf("(%s)", res)
			}
			return res
		}
	case TArr:
		return fmt.Sprintf("(%s -> %s)", ShowType(tT.T1), ShowType(tT.T2))
	case TApp:
		return fmt.Sprintf("(%s %s)", ShowType(tT.T1), ShowType(tT.T2))
	case TTpl:
		{
			res := ""

			for _, t := range tT.Ts {
				if res != "" {
					res = fmt.Sprintf("%s, %s", res, ShowType(t))
				} else {
					res = ShowType(t)
				}
			}

			return fmt.Sprintf("(%s)", res)
		}
	default:
		// Cannot happen because we've covered all `Type` ADT
		// constructors.
		panic("impossible")
	}
}
