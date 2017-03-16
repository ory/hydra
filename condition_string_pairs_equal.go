
package ladon

// StringPairsEqualCondition is a condition which is fulfilled if the given
// array of pairs contains two-element string arrays where both elements
// in the string array are equal
type StringPairsEqualCondition struct{}

// Fulfills returns true if the given value is an array of string arrays and
// each string array has exactly two values which are equal
func (c *StringPairsEqualCondition) Fulfills(value interface{}, _ *Request) bool {
  pairs, PairsOk := value.([]interface{})

  if PairsOk {
    for _, v := range pairs {
      pair, PairOk := v.([]interface{})
      if !PairOk || (len(pair) != 2) {
        return false
      }

      a, AOk := pair[0].(string)
      b, BOk := pair[1].(string)

      if !AOk || !BOk || (a != b) {
        return false
      }
    }
    return true
  }

  return false
}

// GetName returns the condition's name.
func (c *StringPairsEqualCondition) GetName() string {
  return "StringPairsEqualCondition"
}
