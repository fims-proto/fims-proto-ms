package domain

func (v *Voucher) Post() error {
	if v.IsPosted() {
		return newDomainErr(errPostRepeatPost)
	}
	if !v.IsAudited() {
		return newDomainErr(errPostNotAudited)
	}
	if !v.IsReviewed() {
		return newDomainErr(errPostNotReviewed)
	}
	v.isPosted = true
	return nil
}
