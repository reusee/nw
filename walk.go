package hns

type WalkFunc func(*WalkCtx, *Node)

type WalkCtx struct {
}

func (n *Node) Walk(fn WalkFunc) {
	ctx := &WalkCtx{}
	fn(ctx, n)
}

func Do(fn func(*Node)) WalkFunc {
	return func(ctx *WalkCtx, node *Node) {
		fn(node)
	}
}

func Descendant(predict func(*Node) bool, cont WalkFunc) WalkFunc {
	var f func(*WalkCtx, *Node)
	f = func(ctx *WalkCtx, node *Node) {
		for _, child := range node.Children {
			if predict(child) {
				cont(ctx, child)
			} else {
				f(ctx, child)
			}
		}
	}
	return f
}

func DescendantTagEq(tag string, cont WalkFunc) WalkFunc {
	return Descendant(func(node *Node) bool {
		return node.Tag == tag
	}, cont)
}