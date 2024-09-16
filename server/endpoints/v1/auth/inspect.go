package auth

func (h *Auth) handleInspect(restrict bool) (string, any) {
	return "successfully get server authentication status", map[string]any{
		"restrict": restrict,
	}
}
