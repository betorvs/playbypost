interface Validator {
    id: number;
    kind: string;
    valid: boolean;
    output: string;
    checksum: string;
    updated_at: string;
    analise: Analise;
}

interface Analise {
    results: string[];
}

export default Validator;