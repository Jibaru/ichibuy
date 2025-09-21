package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"

	"ichibuy/store/internal/services"
)

// GraphQLStores godoc
// @Summary      GraphQL endpoint for stores
// @Description  GraphQL endpoint to query stores with filters, sorting and pagination
// @Tags         graphql
// @Accept       json
// @Produce      json
// @Param        query body object{query=string,variables=object} true "GraphQL query"
// @Success      200  {object}  object
// @Failure      400  {object}  ErrorResp
// @Failure      401  {object}  ErrorResp
// @Router       /api/v1/graphql [post]
// @Security     BearerAuth
func GraphQLStores(listStoresService *services.ListStores) gin.HandlerFunc {
	locationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Location",
		Fields: graphql.Fields{
			"lat": &graphql.Field{Type: graphql.Float},
			"lng": &graphql.Field{Type: graphql.Float},
		},
	})

	storeType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Store",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.String},
			"name":        &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"location":    &graphql.Field{Type: locationType},
			"slug":        &graphql.Field{Type: graphql.String},
			"createdAt":   &graphql.Field{Type: graphql.String},
			"updatedAt":   &graphql.Field{Type: graphql.String},
		},
	})

	storeListType := graphql.NewObject(graphql.ObjectConfig{
		Name: "StoreList",
		Fields: graphql.Fields{
			"stores": &graphql.Field{Type: graphql.NewList(storeType)},
			"total":  &graphql.Field{Type: graphql.Int},
			"offset": &graphql.Field{Type: graphql.Int},
			"limit":  &graphql.Field{Type: graphql.Int},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"stores": &graphql.Field{
				Type: storeListType,
				Args: graphql.FieldConfigArgument{
					"name":        &graphql.ArgumentConfig{Type: graphql.String},
					"description": &graphql.ArgumentConfig{Type: graphql.String},
					"sortBy":      &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "name"},
					"sortOrder":   &graphql.ArgumentConfig{Type: graphql.String, DefaultValue: "ASC"},
					"offset":      &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 0},
					"limit":       &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 10},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					ctx := params.Context

					userID := ctx.Value("user_id")
					if userID == nil {
						return nil, fmt.Errorf("user not authenticated")
					}

					filters := services.StoreFilters{
						UserID: userID.(string),
					}

					if name, ok := params.Args["name"].(string); ok && name != "" {
						filters.Name = &name
					}

					if description, ok := params.Args["description"].(string); ok && description != "" {
						filters.Description = &description
					}

					offset := params.Args["offset"].(int)
					limit := params.Args["limit"].(int)

					pagination := services.Pagination{
						Offset: offset,
						Limit:  limit,
					}

					sortBy := params.Args["sortBy"].(string)
					sortOrder := params.Args["sortOrder"].(string)

					sorting := services.Sorting{
						Field: sortBy,
						Order: sortOrder,
					}

					serviceReq := services.ListStoresReq{
						Filters:    filters,
						Pagination: pagination,
						Sorting:    sorting,
					}

					resp, err := listStoresService.Exec(ctx, serviceReq)
					if err != nil {
						return nil, err
					}

					return resp, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		var requestBody struct {
			Query     string                 `json:"query"`
			Variables map[string]interface{} `json:"variables"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResp{Error: err.Error()})
			return
		}

		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, ErrorResp{Error: "user not found in context"})
			return
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "user_id", userID)

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  requestBody.Query,
			VariableValues: requestBody.Variables,
			Context:        ctx,
		})

		c.JSON(http.StatusOK, result)
	}
}
