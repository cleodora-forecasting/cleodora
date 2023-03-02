import {FC} from "react";
import {Chip} from "@mui/material";
import {Resolution} from "./__generated__/graphql";


export const ResolutionChip: FC<{resolution: Resolution}> = ({resolution}) => {
    let background: string;
    let label: string;
    label = resolution;
    switch (resolution) {
        case Resolution.Resolved:
            background = "darkgreen";
            break
        case Resolution.Unresolved:
            background = "darkred";
            break
        case Resolution.NotApplicable:
            background = "gray";
            label = "N/A";
            break
        default:
            background = "gray";
            break
    }
    return <Chip label={label} style={{backgroundColor: background, color: "white"}} size="medium" />
}
