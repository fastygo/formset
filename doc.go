// Package formset is the schema kernel for typed admin forms.
//
// A record type plus locale documents (payload_ru / payload_en) bind into a
// Form that a slot can render. Saving writes the documents back as JSON.
//
// This module has no renderer, BFF, or storage dependency.
package formset
