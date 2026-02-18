package account

func (a *Account) UpdateLeaf(isLeaf bool) error {
	a.isLeaf = isLeaf
	return nil
}
