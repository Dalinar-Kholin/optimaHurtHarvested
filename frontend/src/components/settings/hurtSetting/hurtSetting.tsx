import {Accordion, AccordionDetails, AccordionSummary, Snackbar} from "@mui/material";
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {hurtNames, hurtNamesIterable} from "../../../interfaces.ts";
import HurtComp from "./hurtComp.tsx";
import {useState} from "react";

interface IHurtSettings{
    fn : (username: string, pass : string, name : hurtNames) => void
}


export default function HurtSetting({fn} : IHurtSettings) {

    const availableHurtGetResult = localStorage.getItem("availableHurts")
    const availableHurt = availableHurtGetResult!==null ? +availableHurtGetResult : 0 // nie tylkać bo kompilator spadnie z rowerka

    const [state, setState] = useState<boolean>(false);
    const [message, setMessage] = useState<string>("")




    return (
        <div>
            <p>Dostępne hurtownie</p>
            {hurtNamesIterable.map(name =>
                {

                    return name === hurtNames.none ? <></> : <Accordion>
                    <AccordionSummary
                        expandIcon={<ExpandMoreIcon/>}
                        aria-controls="panel1-content"
                        id="panel1-header"
                        sx={{backgroundColor:(availableHurt&name) > 0 ? "#81c784" : "" }}
                    >
                        {hurtNames[name] + ((availableHurt&name)>0 ? " - obecne dane są poprawne" : " - brak poprawnych danych")}
                    </AccordionSummary>
                    <AccordionDetails>
                        <HurtComp fn={(username, pass, name) => {
                            setState(true)
                            setMessage("hasło od " + hurtNames[name] + " dodano do zapisania")
                            fn(username, pass, name)

                        }} name={name}
                        />
                    </AccordionDetails>
                </Accordion>}
                )
            }
            <Snackbar
                anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
                open={state}
                onClose={()=> setState(false)}
                message={message}
            />

        </div>
    )
}