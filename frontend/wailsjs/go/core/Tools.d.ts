// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {structs} from '../models';

export function AntivirusIdentify(arg1:string):Promise<Array<structs.AntivirusResult>>;

export function ConnectAndExecute(arg1:string,arg2:string,arg3:string,arg4:string,arg5:string):Promise<string>;

export function ExtractAlibabaDruidWebSession(arg1:string):Promise<Array<string>>;

export function ExtractAlibabaDruidWebURI(arg1:string):Promise<Array<string>>;

export function ExtractIP(arg1:string):Promise<string>;

export function ExtractURLs(arg1:string):Promise<Array<string>>;

export function FormatOutput(arg1:string):Promise<{[key: string]: Array<string>}>;

export function GOOS():Promise<string>;

export function GetToken(arg1:string,arg2:string,arg3:string):Promise<string>;

export function IPParse(arg1:Array<string>):Promise<Array<string>>;

export function IsRoot():Promise<boolean>;

export function PatchIdentify(arg1:string):Promise<Array<structs.AuthPatch>>;

export function PortParse(arg1:string):Promise<Array<number>>;
