import { TextField} from "@mui/material";
import Box from "@mui/material/Box";
import { styled } from '@mui/system';

interface IHurtResultForm {
    isCheaper: boolean,
    name: string,
    princeForOne: number,
    priceForPack: number,
    productsInPack: number
}


const CheaperComp = styled(TextField)(({}) => ({
    "& .MuiOutlinedInput-root.Mui-disabled": {
        "& fieldset": {
            borderColor: "green",
        },
    },
    "& .MuiInputLabel-root.Mui-disabled": {
        color: "green",
    },
    "& .MuiInputBase-root.Mui-disabled": {
        color: "green",
    },
}));


export default function HurtResultForm({name, princeForOne, priceForPack, productsInPack, isCheaper}: IHurtResultForm) {
    return (
        <>
            {isCheaper?<Box sx={{display: "flex", gap: "10px", justifyContent: "space-around"}} >
                <CheaperComp
                    disabled
                    id="outlined-disabled"
                    label="nazwa Hurtowni"
                    value={name}
                />
                <CheaperComp
                    disabled
                    id="outlined-disabled"
                    label="cena za sztukę"
                    value={princeForOne === -1 ? "Brak produktu" : princeForOne}
                />
                <CheaperComp
                    disabled
                    id="outlined-disabled"
                    label={"Cena za " + (productsInPack === -1 ? "" : productsInPack)}
                    value={priceForPack === -1 ? "Brak produktu" : priceForPack.toFixed(2)}
                />
            </Box>
                :<Box sx={{display: "flex", gap: "10px", justifyContent: "space-around"}} >
                <TextField
                    disabled
                    id="outlined-disabled"
                    label="nazwa Hurtowni"
                    value={name}
                />
                <TextField
                    disabled
                    id="outlined-disabled"
                    label="cena za sztukę"
                    value={princeForOne === -1 ? "Brak produktu" : princeForOne}
                />
                <TextField
                    disabled
                    id="outlined-disabled"
                    label={"Cena za " + (productsInPack === -1 ? "" : productsInPack)}
                    value={priceForPack === -1 ? "Brak produktu" : priceForPack.toFixed(2)}
                />
            </Box> }

        </>
    )
}