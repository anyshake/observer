package helicorder

func (s *HelicorderService) loadDefault() {
	s.imageSize = HELICORDER_IMAGE_SIZE
	s.spanSamples = HELICORDER_DOWNSAMPLE_FACTOR
	s.scaleFactor = HELICORDER_SCALE_FACTOR
}
