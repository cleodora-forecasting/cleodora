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
    const {data} = useQuery(GET_METADATA);
    const [config, setConfig] = useState({} as FrontendConfig);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = () => {
        fetch("config.json")
            .then(res => res.json())
            .then(
                (result) => {
                    setConfig(result as FrontendConfig)
                 }
            )
            .catch(
                (reason) => {
                    console.log("Error getting config", reason)
                }
            )
    };

    return <footer>
        <Grid container spacing={2}>
            <Grid item lg={6}>
                Cleodora Forecasting (<a href="https://cleodora.org">cleodora.org</a>). Version: {data?.metadata.version}
            </Grid>
            <Grid item lg={6}>
                <Box display="flex" justifyContent="flex-end">
                    {config.FOOTER_TEXT}
                </Box>
            </Grid>
        </Grid>
    </footer>
}
