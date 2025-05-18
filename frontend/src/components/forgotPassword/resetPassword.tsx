import {Alert, Button, TextField} from "@mui/material";
import {useState} from "react";
import {useLocation} from "react-router-dom";

export default function ResetPassword() {
    const [password, setPassword] = useState<string>("")
    const [confirmPassword, setConfirmPassword] = useState<string>("")
    const [message, setMessage] = useState<string>("")
    const location = useLocation();

    // Tworzymy instancję URLSearchParams na podstawie location.search
    const queryParams = new URLSearchParams(location.search);

    // Wyciągamy wartość tokenu
    const token = queryParams.get('token');
    const sentData = async () => {
        if (password!== confirmPassword || token===""){
            return
        }

        const body = {
            password: password,
            token: token,
        }
        fetch("/api/resetPassword",{
            method: "POST",
            body : JSON.stringify(body)
        }).then(res => {
            return res.json()
        }). then(data => {

            if (data.error != undefined){
                setMessage(data.error)
                return
            }

            setMessage(data.message)
        }).catch(() => {
            setMessage("nie udało się zmienić hasła")
        })


    }


    return (
        <>
            <h1>wprowadź nowe hasło</h1>
            <TextField label={"hasło"} value={password} onChange={e => setPassword(e.target.value)} autoComplete="off"/>
            <TextField label={"hasło ponownie"} value={confirmPassword} onChange={e => setConfirmPassword(e.target.value)} autoComplete="off"/>
            <Button onClick={sentData}>potwierdź</Button>
            {message!== "" ? <Alert>{message}</Alert> : <></> }

        </>

    )
}