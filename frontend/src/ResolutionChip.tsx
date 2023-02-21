import {FC} from "react";
import {Chip} from "@mui/material";
import {Resolution} from "./__generated__/graphql";


export const ResolutionChip: FC<{resolution: Resolution}> = ({resolution}) => {
    let background: string;
    switch (resolution) {
        case Resolution.Resolved:
            background = "darkgreen";
            break
        case Resolution.Unresolved:
            background = "darkred";
            break
        default:
            background = "gray";
            break
    }
    return <Chip label={resolution} style={{backgroundColor: background, color: "white"}} size="medium" />
}
