package adm

type admQueryOpts struct {
	lv              *int
	startAfterFid   *string
	startAfterId    *string
	batchSize       int
	includeGeometry bool
}

type admQueryOptsBuilder struct {
	conf admQueryOpts
}

func NewAdmQueryOptsBuilder() *admQueryOptsBuilder {
	return &admQueryOptsBuilder{conf: admQueryOpts{
		batchSize:       100,
		includeGeometry: false,
	}}
}

func (builder *admQueryOptsBuilder) SetLv(lv int) *admQueryOptsBuilder {
	builder.conf.lv = &lv
	return builder
}

func (builder *admQueryOptsBuilder) SetStartAfterId(startAfterId string) *admQueryOptsBuilder {
	builder.conf.startAfterId = &startAfterId
	return builder
}

func (builder *admQueryOptsBuilder) SetStartAfterFid(startAfterFid string) *admQueryOptsBuilder {
	builder.conf.startAfterFid = &startAfterFid
	return builder
}

func (builder *admQueryOptsBuilder) SetBatchSize(batchSize int) *admQueryOptsBuilder {
	builder.conf.batchSize = batchSize
	return builder
}

func (builder *admQueryOptsBuilder) Build() (admQueryOpts, error) {
	return builder.conf, nil
}
