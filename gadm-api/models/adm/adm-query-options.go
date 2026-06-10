package adm

import "gadm-api/utils"

type admQueryOpts struct {
	lv              *int
	startAfterFid   *string
	startAfterId    *string
	batchSize       *int
	includeGeometry bool
}

type admQueryOptsBuilder struct {
	conf admQueryOpts
}

func NewAdmQueryOptsBuilder() *admQueryOptsBuilder {
	batchSize := 100
	return &admQueryOptsBuilder{conf: admQueryOpts{
		batchSize:       &batchSize,
		includeGeometry: false,
	}}
}

func (builder *admQueryOptsBuilder) SetStartAfterId(startAfterId string) *admQueryOptsBuilder {
	if startAfterId == "" {
		return builder
	}

	builder.conf.startAfterId = &startAfterId
	return builder
}

func (builder *admQueryOptsBuilder) SetStartAfterFid(startAfterFid string) *admQueryOptsBuilder {
	if startAfterFid == "" {
		return builder
	}
	builder.conf.startAfterFid = &startAfterFid
	return builder
}

func (builder *admQueryOptsBuilder) SetLvAndBatchSize(lv *int, batchSize *int) *admQueryOptsBuilder {
	builder.conf.lv = lv

	_batchSize := *batchSize
	if lv != nil {
		if *lv < 2 {
			_batchSize = utils.Clamp(_batchSize, 1, 5)
		} else if *lv < 4 {
			_batchSize = utils.Clamp(_batchSize, 1, 20)
		} else {
			_batchSize = utils.Clamp(_batchSize, 1, 50)
		}
	}

	builder.conf.batchSize = &_batchSize
	return builder
}

func (builder *admQueryOptsBuilder) Build() (admQueryOpts, error) {
	return builder.conf, nil
}

func (builder *admQueryOptsBuilder) SetIncludeGeometry(includeGeometry bool) *admQueryOptsBuilder {
	builder.conf.includeGeometry = includeGeometry
	return builder
}
