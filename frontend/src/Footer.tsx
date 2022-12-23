import {FC} from "react";
import {gql} from "./__generated__"
import "./Footer.css"
import {useQuery} from "@apollo/client";
import {Box, Grid} from "@mui/material";

export const GET_METADATA = gql(`
    query GetMetadata {
        metadata {
            version
        }
    }
`);

export const Footer: FC = () => {
    const {data} = useQuery(GET_METADATA);

    return <footer>
        <Grid container spacing={2}>
            <Grid item lg={6}>
                Cleodora Forecasting (<a href="https://cleodora.org">cleodora.org</a>). Version: {data?.metadata.version}
            </Grid>
            <Grid item lg={6}>
                <Box display="flex" justifyContent="flex-end">
                </Box>
            </Grid>
        </Grid>
    </footer>
}
