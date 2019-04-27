package parse

import (
	"testing"
	"strings"
)

func TestReadXml(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<types xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="/home/johan/rti_connext_dds-5.2.3/bin/../resource/app/app_support/rtiddsgen/schema/rti_dds_topic_types.xsd">
<module name="com">
  <module name="jhbe">
    <module name="example">
      <const name="A" type="string" value="&quot;A String&quot;"/>
      <const name="B" type="long" value="5"/>
      <enum name="C">
        <enumerator name="C_One"/>
        <enumerator name="C_Two"/>
      </enum> 
      <const name="DD" type="nonBasic"  nonBasicTypeName= "com::jhbe::example::C" value="(com::jhbe::example::C_Two)"/>
      <struct name= "D">
        <member name="D_A" type="long"/>
        <member name="D_B" type="float"/>
        <member name="D_C" type="string"/>
        <member name="D_D" type="boolean"/>
        <member name="D_E" sequenceMaxLength="10" type="boolean"/>
      </struct>
      <union name="E">
        <discriminator type="nonBasic" nonBasicTypeName="com::jhbe::example::C"/>
        <case>
          <caseDiscriminator value="(com::jhbe::example::C_One)"/>
          <member name="E_A" type="nonBasic"  nonBasicTypeName= "com::jhbe::example::D"/>
        </case>
        <case>
          <caseDiscriminator value="(com::jhbe::example::C_Two)"/>
          <member name="E_B" type="nonBasic"  nonBasicTypeName= "com::jhbe::example::C"/>
        </case>
      </union>
      <struct name= "F" baseType="com::jhbe::example::D">
        <member name="F_A" type="short"/>
        <member name="F_B" type="unsignedShort"/>
        <member name="F_C" type="long"/>
        <member name="F_D" type="unsignedLong"/>
        <member name="F_E" type="float"/>
        <member name="F_F" type="double"/>
        <member name="F_G" type="boolean"/>
        <member name="F_H" type="string"/>
        <member name="F_I" type="nonBasic"  nonBasicTypeName= "com::jhbe::example::C"/>
        <member name="F_J" type="nonBasic"  nonBasicTypeName= "com::jhbe::example::D"/>
        <member name="F_K" sequenceMaxLength="com::jhbe::example::B" type="nonBasic"  nonBasicTypeName= "com::jhbe::example::C"/>
        <member name="F_L" type="nonBasic"  nonBasicTypeName= "com::jhbe::example::E"/>
      </struct>
    </module>
  </module>
</module>
<module name="this">
  <module name="that">
    <const name="dummy" type="long" value="2"/>
  </module>
</module>
</types>
`

	types, err := ReadXml(strings.NewReader(xml))
	if err != nil {
		t.Error("Did not expect an error.")
	}
	if len(types.ModuleElements.Modules[0].Modules[0].Modules[0].Consts) != 3 {
		t.Error("Expected com::jhbe::types to have three consts.")
	}
}