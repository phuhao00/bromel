package meshwork

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/phuhao00/bromel/logger"
	"os"
	"reflect"
	"runtime"
)

func newRouter() *Router {
	return &Router{
		stopped:           make(chan struct{}),
		handlerMapper:     make(map[interface{}]HandlerFunc),
		middlewaresMapper: make(map[interface{}][]MiddlewareFunc),
	}
}

// Router is a router for incoming message.
// Router routes the message to its handler and middlewares.
type Router struct {
	// handlerMapper maps message's ID to handler.
	// Handler will be called around middlewares.
	handlerMapper map[interface{}]HandlerFunc

	// middlewaresMapper maps message's ID to a list of middlewares.
	// These middlewares will be called before the handler in handlerMapper.
	middlewaresMapper map[interface{}][]MiddlewareFunc

	// globalMiddlewares is a list of MiddlewareFunc.
	// globalMiddlewares will be called before the ones in middlewaresMapper.
	globalMiddlewares []MiddlewareFunc

	notFoundHandler HandlerFunc
	reqCtxQueue     chan *Context
	stopped         chan struct{}
}

// HandlerFunc is the function type for handlers.
type HandlerFunc func(ctx *Context) (*Entry, error)

// MiddlewareFunc is the function type for middlewares.
// A common pattern is like:
//
// 	var md MiddlewareFunc = func(next HandlerFunc) HandlerFunc {
// 		return func(ctx *Context) (message.Entry, error) {
// 			return next(ctx)
// 		}
// 	}
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

var nilHandler HandlerFunc = func(ctx *Context) (*Entry, error) {
	return nil, nil
}

func (r *Router) stop() {
	close(r.stopped)
}

func (r *Router) consumeRequest() {
	defer fmt.Printf("router stopped")
	for {
		select {
		case <-r.stopped:
			close(r.reqCtxQueue)
			return
		case reqCtx, ok := <-r.reqCtxQueue:
			if !ok {
				return
			}
			select {
			case <-reqCtx.session.closed:
				reqCtx.session.ctxPool.Put(reqCtx)
				continue
			default:
			}
			if reqCtx.reqMsgEntry == nil {
				reqCtx.session.ctxPool.Put(reqCtx)
				continue
			}

			go func() {
				defer reqCtx.session.ctxPool.Put(reqCtx)
				respEntry, err := r.handleRequest(reqCtx)
				if err != nil {
					fmt.Printf("router handle request err: %s", err)
					return
				}
				if respEntry == nil {
					return
				}
				if err := reqCtx.session.SendResp(respEntry); err != nil {
					fmt.Printf("router send resp err: %s", err)
				}
			}()
		}
	}
}

func (r *Router) handleRequest(ctx *Context) (*Entry, error) {
	var handler HandlerFunc
	if v, has := r.handlerMapper[ctx.reqMsgEntry.ID]; has {
		handler = v
	}

	var mws = r.globalMiddlewares
	if v, has := r.middlewaresMapper[ctx.reqMsgEntry.ID]; has {
		mws = append(mws, v...) // append to global ones
	}

	// create the handlers stack
	wrapped := r.wrapHandlers(handler, mws)

	// and call the handlers stack
	return wrapped(ctx)
}

// wrapHandlers wraps handler and middlewares into a right order call stack.
// Makes something like:
// 	var wrapped HandlerFunc = m1(m2(m3(handle)))
func (r *Router) wrapHandlers(handler HandlerFunc, middles []MiddlewareFunc) (wrapped HandlerFunc) {
	if handler == nil {
		handler = r.notFoundHandler
	}
	if handler == nil {
		handler = nilHandler
	}
	wrapped = handler
	for i := len(middles) - 1; i >= 0; i-- {
		m := middles[i]
		wrapped = m(wrapped)
	}
	return wrapped
}

// register stores handler and middlewares for id.
func (r *Router) register(id interface{}, h HandlerFunc, m ...MiddlewareFunc) {
	if h != nil {
		r.handlerMapper[id] = h
	}
	ms := make([]MiddlewareFunc, 0, len(m))
	for _, mm := range m {
		if mm != nil {
			ms = append(ms, mm)
		}
	}
	if len(ms) != 0 {
		r.middlewaresMapper[id] = ms
	}
}

// registerMiddleware stores the global middlewares.
func (r *Router) registerMiddleware(m ...MiddlewareFunc) {
	for _, mm := range m {
		if mm != nil {
			r.globalMiddlewares = append(r.globalMiddlewares, mm)
		}
	}
}

// printHandlers prints registered route handlers to console.
func (r *Router) printHandlers(addr string) {
	fmt.Printf("\n[PATIENCE DESERT ROUTE TABLE]:\n")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Message ID", "Route Handler"})
	table.SetAutoFormatHeaders(false)
	for id, h := range r.handlerMapper {
		handlerName := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		table.Append([]string{fmt.Sprintf("%v", id), handlerName})
	}
	table.Render()
	fmt.Printf("[PATIENCE DESERT ROUTE TABLE] Serving at: %s\n\n", addr)
}

func (r *Router) setNotFoundHandler(handler HandlerFunc) {
	r.notFoundHandler = handler
}

func RecoverMiddleware(log *logger.Logger) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) (*Entry, error) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			return next(c)
		}
	}
}
