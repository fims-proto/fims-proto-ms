package account

func (a *Account) UpdateNumber(levelNumber int, codeLengths []int) error {
	a.numberHierarchy = a.numberHierarchy[:len(a.numberHierarchy)-1]
	a.numberHierarchy = append(a.numberHierarchy, levelNumber)

	accountNumber, err := composeAccountNumber(a.numberHierarchy, codeLengths)
	if err != nil {
		return err
	}

	a.accountNumber = accountNumber
	return nil
}
