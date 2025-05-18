import LoginForm from "./components/login/loginForm.tsx";
import { Navigate, Route, Routes, useLocation} from 'react-router-dom';
import MainSite from "./components/home/mainSite.tsx";
import AppBarCustomed from "./components/appBar/appBar.tsx";
import {createTheme, CssBaseline, ThemeProvider} from "@mui/material";
import useCheckCookie from "./customHook/useCheckCookie.ts";
import Tariff from "./components/tariff/Tariff.tsx";
import SettingPage from "./components/settings/settings.tsx";
import {useEffect, useState} from "react";
import FailedLoginPage from "./components/failedLoginPage/failedLoginPage.tsx";
import SignIn from "./components/SignIn.tsx";
import PaymentComp from "./components/payment/paymentComp.tsx";
import ForgotPassword from "./components/forgotPassword/forgotPassword.tsx";
import {freeBarAndCookiePath} from "./interfaces.ts";
import ResetPassword from "./components/forgotPassword/resetPassword.tsx";
import fetchWithAuth from "./typeScriptFunc/fetchWithAuth.ts";
import ContactComp from "./components/contact/contact.tsx";


const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});


function CheckCookie() {
    useCheckCookie();
    return <></>
}

const sendRequest = async () => {
    try {
        const response = await fetch('/api/isAlive', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        });
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
    } catch (error) {
        console.error('Error during fetch:', error);
    }
};

const logout = () => {
    fetchWithAuth("/api/logout").then(
        ()=>{
            console.log("logouted")
        }
    )
}

function App() {





    const [showAppBar, setShowAppBar] = useState<boolean>(false)

    const location = useLocation();

    useEffect(() => {
        if (freeBarAndCookiePath.some(i => i === location.pathname)) {
            setShowAppBar(false);
        } else {
            setShowAppBar(true);
        }
    }, [location.pathname]);

    useEffect(() => {
        const intervalId = setInterval(sendRequest, 240000);

        // Wyczyszczenie interwału przy odmontowywaniu komponentu, aby nie występowało wielokrotne wysyłanie
        return () => clearInterval(intervalId);
    }, []);

    useEffect(() => {
        window.addEventListener('beforeunload', logout);

        return () => {
            window.removeEventListener('beforeunload', logout);
        };
    }, []);

    return (
        <div>
            <ThemeProvider theme={darkTheme}>
                <CssBaseline/>
                {/*<Router> komponent App jest owrapowany więc tutaj nie trzeba*/}
                    <CheckCookie/>
                    {showAppBar ?<AppBarCustomed /> : <></> }
                    <Routes>
                        <Route path={"/login"} element={<LoginForm/>}/>
                        <Route path={"/strona główna"} element={<MainSite/>}/>
                        <Route path={"/cennik"} element={<Tariff/>}/>
                        <Route path={"/ustawienia"} element={<SettingPage/>}/>
                        <Route path={"/*"} element={<Navigate to={"/strona główna"}/>}/>
                        <Route path={"/failed"} element={<FailedLoginPage/>}/>
                        <Route path={"/signIn"} element={<SignIn/>} />
                        <Route path={"/płatności"} element={<PaymentComp/>}/>
                        <Route path={"/forgotPassword"} element={<ForgotPassword/>}/>
                        <Route path={"/resetPassword"} element={<ResetPassword/>}/>
                        <Route path={"/kontakt"} element={<ContactComp/>}/>
                    </Routes>
            </ThemeProvider>
        </div>
    )
}

export default App
