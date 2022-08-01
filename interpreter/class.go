package interpreter

type LoxClass struct { // impl LoxCallable
	className string
	methods   map[string]*LoxCustomFunc
}

func NewLoxClass(name string, methods map[string]*LoxCustomFunc) *LoxClass {
	return &LoxClass{
		className: name,
		methods:   methods,
	}
}

func (c *LoxClass) String() string {
	return c.className
}

func (c *LoxClass) arity() int {
	initializer := c.findMethod("init")
	if initializer == nil {
		return 0
	}
	return initializer.arity()
}

func (c *LoxClass) call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	instance := NewLoxInstance(c)
	initializer := c.findMethod("init")
	if initializer != nil {
		initializer.bind(instance).call(interpreter, arguments)
	}

	return NewLoxInstance(c), nil
}

func (c *LoxClass) Name() string {
	return c.className
}

func (c *LoxClass) findMethod(name string) (result *LoxCustomFunc) {
	var ok bool
	if result, ok = c.methods[name]; ok {
		return
	}
	return nil
}
