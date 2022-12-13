import {FC} from "react";
import {gql} from "./__generated__"
import "./Footer.css"
import {useQuery} from "@apollo/client";

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
        Version: {data?.metadata.version}
    </footer>
}
