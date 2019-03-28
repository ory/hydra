package x

import (
	"github.com/julienschmidt/httprouter"

	"github.com/ory/x/serverx"
)

type RouterAdmin struct {
	*httprouter.Router
}

type RouterPublic struct {
	*httprouter.Router
}

func (r *RouterPublic) RouterAdmin() *RouterAdmin {
	return &RouterAdmin{Router: r.Router}
}

func (r *RouterAdmin) RouterPublic() *RouterPublic {
	return &RouterPublic{Router: r.Router}
}

func NewRouterPublic() *RouterPublic {
	router := httprouter.New()
	router.NotFound = serverx.DefaultNotFoundHandler
	return &RouterPublic{
		Router: router,
	}
}

func NewRouterAdmin() *RouterAdmin {
	router := httprouter.New()
	router.NotFound = serverx.DefaultNotFoundHandler
	return &RouterAdmin{
		Router: router,
	}
}
