package handler

type Middleware func(next ContextHandler) ContextHandler
