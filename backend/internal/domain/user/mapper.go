package user

type PermissionMapper interface {
	// NOTE: IDP側のRoleが増えるたびにエラー吐くようになると嫌なので
	// error は返さない。
	// マッピングできないものはスルーするような実装が望ましい
	MapToPermissions(roleIDs []string) []Permission
}
