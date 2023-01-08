package golang

func DefaultValue(id Ident) string {
	switch id.From {
	case "":
		switch id.Name {
		case "bool":
			return "false"
		case "string":
			return "\"\""
		case "int":
			return "0"
		}
	}
	return "{}"
}
