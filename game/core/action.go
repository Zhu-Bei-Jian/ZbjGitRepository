package core

type IActionData interface {
	PostActStream(...func()) IActionData
	IsStop() bool
	Stop()
	Front() (func(), bool)
}

type ActionDataCore struct {
	parent            IActionData
	actonStream       []func()
	delayActionStream []func()
	stop              bool
}

func (ad *ActionDataCore) PostActStream(funs ...func()) IActionData {
	ad.delayActionStream = append(ad.delayActionStream, funs...)
	return ad
}

func (ad *ActionDataCore) Stop() {
	ad.stop = true
}

func (ad *ActionDataCore) IsStop() bool {
	return ad.stop
}

func (ad *ActionDataCore) Front() (func(), bool) {
	if len(ad.delayActionStream) > 0 {
		ad.actonStream = append(ad.delayActionStream, ad.actonStream...)
		ad.delayActionStream = nil
	}

	if len(ad.actonStream) == 0 {
		return nil, false
	}

	ret := ad.actonStream[0]
	ad.actonStream = ad.actonStream[1:]
	return ret, true
}

func (ad *ActionDataCore) Parent() IActionData {
	return ad.parent
}

func (ad *ActionDataCore) SetParent(adi IActionData) {
	ad.parent = adi
}
