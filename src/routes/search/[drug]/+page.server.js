import meds from "$lib/data/meds.json";

export function load({ params }) {
    
	return {
        name: params.drug,
		med: meds[params.drug]
	};
}