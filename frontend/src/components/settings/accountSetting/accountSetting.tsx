import {useState} from "react";
import {Alert, AlertTitle, Button, TextField} from "@mui/material";

interface IAccountSetting{
    fn :(s : string)=>void
}

export default function AccountSetting({fn} : IAccountSetting) {

    const [password, setPassword] = useState<string>("")
    const [passwordAgain, setPasswordAgain] = useState<string>("")
    const [error, setError] = useState<string>("")

    return (
        <div style={{display: "flex", justifyContent: "center"}}>
            <TextField value={password} autoComplete={"off"} label={"nowe hasło"} onChange={(e) => {
                setPassword(e.target.value)
            }}/>
            <TextField value={passwordAgain} autoComplete={"off"} label={"hasło ponownie"} onChange={(e) => {
                setPasswordAgain(e.target.value)
            }}/>
            <Button onClick={()=>{
                if(passwordAgain!= password){
                    setError("hasła nie są takie same")
                    return
                }
                fn(password)
                setPassword("")
                setError("")
                setPasswordAgain("")
            }}>zaktualizuj hasło</Button>
            {error===""? <></> : <Alert severity="error">
                <AlertTitle>Error</AlertTitle>
                {error}
            </Alert>}
        </div>
    )
}