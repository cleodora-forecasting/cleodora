import {graphql, GraphQLHandler, GraphQLRequest, GraphQLVariables} from 'msw'

export const handlers: RequestHandler[] = [
    graphql.query('GetMetadata', (req, res, ctx) => {
        return res(
            ctx.data({
                metadata: {
                    version: "test-test",
                    // typename?
                },
            }),
        )
    }),
]
