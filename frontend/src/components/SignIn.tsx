import {useState} from "react";
import {useNavigate} from "react-router-dom";
import {Alert, AlertTitle, Button, TextField} from "@mui/material";

export default function SignIn(){

    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [repeatedPassword, setRepeatedPassword] = useState<string>("");
    const [companyName, setCompanyName] = useState<string>("")
    const [nip, setNip] = useState<string>("")
    const [street, setStreet] = useState<string>("")
    const [localNumber, setLocalNumber] = useState<string>("")
    const [email, setEmail] = useState<string>("")


    const [isProperData, setIsProperData] = useState<boolean>(true)
    const [errorMessage, setErrorMessage] = useState<string>("")

    const navigate = useNavigate()

    return (
        <>
            <h1>załóż konto</h1>
            <form onSubmit={e => {
                e.preventDefault()
                // logowanie sie
                if (password!==repeatedPassword) {
                    setIsProperData(false)
                    setErrorMessage("hasła nie są takie same")
                    return
                }
                const LoginData = {
                    username : username,
                    password : password,
                    email: email,
                    companyName: companyName,
                    nip: nip,
                    street: street,
                    nr: localNumber
                }
                fetch('/api/signIn', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(LoginData),
                }).then(response => {
                    if (!response.ok) {
                        // Próbujemy pobrać dane błędu jako JSON
                        return response.json().then(errorData => {
                            setErrorMessage(errorData.error || "błąd sieci")
                            setIsProperData(false)
                            throw new Error('Error: ' + errorData.error);
                        });
                    }
                    return response.json();
                }).then(data => {

                    if (data.error != undefined){
                        setErrorMessage(data.error)
                        setIsProperData(false)
                        return
                    }


                    setIsProperData(true)
                    setUsername("")
                    setPassword("")
                    navigate("/login")
                }).catch(error => {
                    console.error('There has been a problem with your fetch operation:', error);
                })

            }}>
                <TextField
                    autoComplete={"off"}
                    id="filled"
                    label="Username"
                    placeholder="username"
                    value={username}
                    onChange={e => setUsername(e.target.value)}
                />
                <p></p>
                <TextField
                    id="outlined-password-input first"
                    label="Password"
                    type="password"
                    autoComplete="current-password"
                    value={ password }
                    onChange={e => setPassword(e.target.value)}
                />
                <p></p>
                <TextField
                    id="outlined-password-input second"
                    label="Password again"
                    type="password"
                    autoComplete="current-password"
                    value={ repeatedPassword }
                    onChange={e => setRepeatedPassword(e.target.value)}
                />
                <p></p>
                <TextField
                    autoComplete={"off"}
                    id="filled"
                    label="email"
                    placeholder="email"
                    value={email}
                    onChange={e => setEmail(e.target.value)}/>
                <p></p>
                <TextField
                    autoComplete={"off"}
                    id="filled"
                    label="nazwa firmy"
                    placeholder="nazwa firmy"
                    value={companyName}
                    onChange={e => setCompanyName(e.target.value)}
                />
                <p></p>
                <TextField
                    autoComplete={"off"}
                    id="filled"
                    label="nip"
                    placeholder="nip"
                    value={nip}
                    onChange={e => setNip(e.target.value)}
                />
                <p></p>
                <TextField
                    autoComplete={"off"}
                    id="filled"
                    label="street"
                    placeholder="street"
                    value={street}
                    onChange={e => setStreet(e.target.value)}
                />
                <p></p>
                <TextField
                autoComplete={"off"}
                id="filled"
                label="numer lokalu"
                placeholder="numer lokalu"
                value={localNumber}
                onChange={e => setLocalNumber(e.target.value)}/>
                <p></p>

                {isProperData ? <div></div> : <Alert severity="error">
                    <AlertTitle>Error</AlertTitle>
                    {errorMessage}
                </Alert>}
                <p></p>
                <Button variant="contained" type={"submit"}>
                    signIn
                </Button>

            </form>
        </>
    )

}