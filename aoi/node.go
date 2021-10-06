package aoi

//Node aoi node
type Node struct {
	owner       EntityIF // 属于哪个AOI单元，这里把代表Entity本身的节点也当作一个R=0的AOI单元
	preX, nextX *Node
	preZ, nextZ *Node
	preY, nextY *Node
}

func (n *Node) nextXEntity() EntityIF {
	return n.nextX.owner
}

func (n *Node) nextZEntity() EntityIF {
	return n.nextZ.owner
}

func (n *Node) preXEntity() EntityIF {
	return n.preX.owner
}

func (n *Node) preZEntity() EntityIF {
	return n.preZ.owner
}

func (n *Node) insertBeforeX(newNode *Node) {
	if n.preX != nil {
		n.preX.nextX = newNode
	}
	newNode.preX = n.preX
	n.preX = newNode
	newNode.nextX = newNode
}

func (n *Node) insertBeforeY(newNode *Node) {
	if n.preY != nil {
		n.preY.nextY = newNode
	}
	newNode.preY = n.preY
	n.preY = newNode
	newNode.nextY = newNode
}

func (n *Node) insertBeforeZ(newNode *Node) {
	if n.preZ != nil {
		n.preZ.nextZ = newNode
	}
	newNode.preZ = n.preZ
	n.preZ = newNode
	newNode.nextZ = newNode
}

func shuffle(entityIF EntityIF) {

	fn := func(tag, cursorTag string) {
		thisPos := entityIF.GetPosition().X
		for {
			var cursorNode EntityIF
			if cursorTag == "pre" {
				if tag == "x" {
					cursorNode = entityIF.GetNode().preXEntity()
				}
				if tag == "z" {
					cursorNode = entityIF.GetNode().preZEntity()
				}
			}
			if cursorTag == "next" {
				if tag == "x" {
					cursorNode = entityIF.GetNode().nextXEntity()
				}
				if tag == "z" {
					cursorNode = entityIF.GetNode().nextZEntity()
				}
			}

			if cursorNode == nil {
				break
			}
			cursorPosX := cursorNode.GetPosition().X
			if thisPos < cursorPosX || (thisPos == cursorPosX && entityIF.order() < cursorNode.order()) {
				if entityIF.isTriggerNode() && !cursorNode.isTriggerNode() {
					entityIF.crossed(cursorNode, true)
				} else if !entityIF.isTriggerNode() && cursorNode.isTriggerNode() {
					cursorNode.crossed(entityIF, false)
				}
				if cursorTag == "pre" {
					entityIF.moveToPrevX()
				}
				if cursorTag == "next" {
					entityIF.moveToNextX()
				}

			} else {
				break
			}
		}
	}

	fn("x", "pre")
	fn("z", "pre")
	fn("x", "next")
	fn("z", "next")
}
