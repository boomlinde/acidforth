package machine

import "errors"

const (
	MAC_NORM = iota
	MAC_NAME
	MAC_COMP
)

func ExpandMacros(tokens []string) ([]string, error) {
	var name string
	state := MAC_NORM
	out := make([]string, 0)
	macros := make(map[string][]string)
	getMac := func(token string) []string {
		mac := macros[token]
		if mac == nil {
			return []string{token}
		} else {
			return mac
		}
	}
	for _, token := range tokens {
		switch {
		case state == MAC_NORM:
			if token == ":" {
				state = MAC_NAME
			} else {
				for _, t := range getMac(token) {
					out = append(out, t)
				}
			}
		case state == MAC_NAME:
			if token == ":" || token == ";" {
				return nil, errors.New("Macro name can't be \":\" or \";\"")
			}
			name = token
			state = MAC_COMP
			macros[token] = make([]string, 0)
		case state == MAC_COMP:
			if token == ":" {
				return nil, errors.New("Macro can't contain \":\"")
			}
			if token == ";" {
				state = MAC_NORM
			} else {
				for _, t := range getMac(token) {
					macros[name] = append(macros[name], t)
				}
			}
		}
	}
	if state != MAC_NORM {
		return nil, errors.New("Unclosed macro detected")
	}
	return out, nil
}
