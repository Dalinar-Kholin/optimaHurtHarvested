import {Button, Snackbar} from "@mui/material";

import {AccountStatus} from "../../interfaces.ts";
import Box from "@mui/material/Box";
import CardComponent from "./CardComponent.tsx";
import fetchWithAuth from "../../typeScriptFunc/fetchWithAuth.ts";
import {useState} from "react";


enum prodName{
    "monthlyDefault" = 0,
    "yearly" = 1
}

export default function PaymentComp() {
    const statusString = localStorage.getItem("accountStatus")
    const status = statusString === null ? 0 : +statusString

    const [openSnackbar, setOpenSnackbar] = useState<boolean>(false)
    const [messageFromBackend, setMessageFromBackend] = useState<string>("")


    const cancelSub = async  ()=> {
        fetchWithAuth("/api/payment/stripe/cancel", {
            method: "GET"
        }).then(response =>{
            return response.json()
        }).then(data =>{

            if (data.error != undefined){
                setMessageFromBackend(data.error)
                setOpenSnackbar(true)
                return
            }

            setMessageFromBackend(data.message)
            setOpenSnackbar(true)
        }).catch(_err => {

        })}



    return (
        <>
            <Box sx={{
                width: '100%',
                typography: 'body1',
                padding: "30px 10px",
                margin: "30px auto",
                borderRadius: "20px",
                backgroundColor: "#363636"
            }}>
                <Box sx = {{display: "flex", justifyContent: "space-evenly"}}>
                    <CardComponent prodName={prodName.yearly} header={"rozpocznij roczną subskrypcje"} name={"subskrypcja roczna"} timePeriod={"365 dni"} description={
                        "subskrypcja zapewniająca dostęp do aplikacji przez 365 dni\n- przy zakupie rocznym obowiązuje zniżka 18%\n- podczas trwania subskrypcji nieograniczona liczba zapytań"}/>
                    <CardComponent prodName={prodName.monthlyDefault} header={status===AccountStatus.New? "ROZPOCZNIJ WERSJE PRÓBNĄ" : "wznów subskrypcje"}  name={"miesięczna"} timePeriod={"30 dni"} description={"subskrypcja miesięczna pozwalająca na dostęp do sprawdzania wyników\n- 200 zapytań pojedyńczych\n- 4 zapytania listowne"}/>
                </Box>

                {status===AccountStatus.Active ? <Button onClick={cancelSub}>zakończ subskrypcje</Button> : <></>}


                <Snackbar
                    anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
                    open={openSnackbar}
                    onClose={()=> {
                        setOpenSnackbar(false)
                        setMessageFromBackend("")
                    }}
                    message={messageFromBackend}
                />
            </Box>



        </>
    )
}
