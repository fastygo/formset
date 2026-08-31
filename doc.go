// Package formset is the schema kernel for typed admin forms.
//
// A record type plus one document per locale bind into a Form a slot can
// render. Saving writes those documents back as JSON. There is no fixed
// payload_ru / payload_en pair.
//
// This module has no renderer, BFF, or storage dependency.
package formset
