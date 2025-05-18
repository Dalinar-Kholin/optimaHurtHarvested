import {Button, Card, CardActions, CardContent, Typography} from "@mui/material";
import fetchWithAuth from "../../typeScriptFunc/fetchWithAuth.ts";
import {loadStripe} from "@stripe/stripe-js";

interface ICardComponent{
    header : string,
    name: string,
    timePeriod: string,
    description: string,
    prodName: number
}


const stripePromise = loadStripe("pk_live_51PmZWL03bfZgIVzMKAMyQ6jAw833lwQbpePt4itZJCZAnQ3ZBFT4gfD7I56DHfWQX8og5i1c7AHqEODqq7Xtz5qJ006nm0AoLj");


export default function CardComponent({name, timePeriod, description, header, prodName}: ICardComponent){

    const handleClick = async () => {
        // Pobranie instancji Stripe
        const stripe = await stripePromise;
        if (stripe == null) {
            return
        }
        // Wywołanie backendu, aby utworzyć sesję
        const response = await fetchWithAuth(`/api/payment/stripe?prodName=${prodName+""}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
        });

        const session = await response.json();

        // Przekierowanie na stronę płatności Stripe
        if ("redirectToCheckout" in stripe) {
            const {error} = await stripe.redirectToCheckout({
                sessionId: session.id,
            });

            if (error) {
                console.error('Error redirecting to checkout:', error);
            }
        }

    };



    return(
        <>
            <Card sx={{ minWidth: 275, flex: 1, marginRight: 1 }}>
                <CardContent>
                    <Typography sx={{ fontSize: 14 }} color="text.secondary" gutterBottom>
                        {header}
                    </Typography>
                    <Typography variant="h5" component="div" sx={{ wordWrap: 'break-word' }}>
                        {name}
                    </Typography>
                    <Typography sx={{ mb: 1.5, wordWrap: 'break-word' }} color="text.secondary">
                        {timePeriod}
                    </Typography>
                    <Typography variant="body2" sx={{ wordWrap: 'break-word', whiteSpace: 'pre-line' }}>
                        {description}
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button size="small" onClick={handleClick}>sprawdź</Button>
                </CardActions>
            </Card>

        </>
    )
}