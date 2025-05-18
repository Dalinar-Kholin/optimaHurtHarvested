import {hurtNames} from "../../../interfaces.ts";

export interface IHurtInfoForComp {
    name: string
    hurtName: hurtNames
    priceForPack: number,
    priceForOne: number,
    productsInPack: number
}
