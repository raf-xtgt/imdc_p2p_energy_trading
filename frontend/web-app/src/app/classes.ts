export class User{
    constructor (public username:string, public email:string, public password:string, public address:string, public smartMeterNo:number){}
}


export class Token{
    constructor (public token:string){}
}