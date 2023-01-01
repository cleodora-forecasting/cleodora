import {
    rest,
    graphql,
    RequestHandler
} from 'msw'

export const handlers: RequestHandler[] = [
    graphql.query('GetMetadata', (req, res, ctx) => {
        console.log("GetMetadata called");
        return res(
            ctx.data({
                metadata: {
                    version: "99.99.99+test",
                    __typename: "Metadata",
                },
            }),
        )
    }),
    rest.get('/config.json', (req, res, ctx) => {
        console.log("config.json called");
        return res(ctx.json({ FOOTER_TEXT: "Footer text for a test" }))
    }),
]
