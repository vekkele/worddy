package icons

type iconProps struct {
	Size   string
	Height string
}

func NewIconProps(size string) iconProps {
	return iconProps{
		Size:   size,
		Height: size,
	}
}

func (p iconProps) WithHeight(height string) iconProps {
	p.Height = height
	return p
}
