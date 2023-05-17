import {FC, useEffect, useState} from "react";
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

type FrontendConfig = {
  FOOTER_TEXT: string;
};

export const Footer: FC = () => {
    const {data, error} = useQuery(GET_METADATA);
    if (error) {
        console.log("Error GET_METADATA", error);
    }
    const [config, setConfig] = useState({} as FrontendConfig);

    useEffect(() => {
        const fetchData = async () => {
            const data = await fetch("config.json");
            const json = await data.json();
            setConfig(json as FrontendConfig);
        };
        fetchData()
            .catch(
                (reason) => {
                    console.error(`Error getting config.json: ${reason}`)
                }
            );
    }, []);


    return <footer>
        <Grid container spacing={2}>
            <Grid item lg={6}>
                Cleodora Forecasting
                (<a href="https://cleodora.org">cleodora.org</a>).
                Version: <span>{data?.metadata.version}</span>
            </Grid>
            <Grid item lg={6}>
                <Box display="flex" justifyContent="flex-end">
                    {config.FOOTER_TEXT}
                </Box>
            </Grid>
        </Grid>
    </footer>
}
