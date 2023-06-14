import {
    rest,
    graphql,
    RequestHandler
} from 'msw'

export const handlers: RequestHandler[] = [
    graphql.query('GetMetadata', async (req, res, ctx) => {
        console.log("GetMetadata called");
        await new Promise(f => setTimeout(f, 200));
        return res(
            ctx.data({
                metadata: {
                    version: "99.99.99+test",
                    __typename: "Metadata",
                },
            }),
        )
    }),
    graphql.query('GetForecasts', async (req, res, ctx) => {
        await new Promise(f => setTimeout(f, 200));
        return res(
            ctx.data({
                "forecasts":[
                    {
                        "id":"1",
                        "title":"Will \"The Fabelmans\" win \"Best Picture\" at the Oscars 2023?",
                        "description":"",
                        "created":"2022-10-30T17:05:00+01:00",
                        "closes":null,
                        "resolves":"2023-03-11T23:59:00+01:00",
                        "resolution":"UNRESOLVED",
                        "estimates":[
                            {
                                "id":"1",
                                "reason": "It's the type of film the Academy" +
                                    " loves",
                                "created": "2022-10-30T17:05:00+01:00",
                                "brierScore": null,
                                "probabilities":[
                                    {
                                        "id":"1",
                                        "value":30,
                                        "outcome":{
                                            "id":"1",
                                            "text":"Yes",
                                            "correct": false,
                                            "__typename":"Outcome"
                                        },
                                        "__typename":"Probability"
                                    },
                                    {
                                        "id":"2",
                                        "value":70,
                                        "outcome":{
                                            "id":"2",
                                            "text":"No",
                                            "correct": false,
                                            "__typename":"Outcome"
                                        },
                                        "__typename":"Probability"
                                    }
                                ],
                                "__typename":"Estimate"
                            }
                        ],
                        "__typename":"Forecast"
                    }
                ]
            }),
        );
    }),
    graphql.mutation("createForecast", async (req, res, ctx) => {
        await new Promise(f => setTimeout(f, 200));
        return res(
            ctx.data({
                "createForecast":{
                    "id": "999",
                    "title": "Mock title",
                    "__typename": "Forecast"
                }
            }),
        )
    }),
    rest.get('/config.json', async (req, res, ctx) => {
        console.log("config.json called");
        await new Promise(f => setTimeout(f, 200));
        return res(ctx.json({ FOOTER_TEXT: "Footer text for a test" }))
    }),
]
