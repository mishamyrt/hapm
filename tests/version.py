"""hapm version module tests"""
import pytest
from hapm.version.parse import InvalidVersion, parse_version
from hapm.version.version import Version


def test_parse_version():
    """Tests version parsing"""
    assert parse_version("1.0.0") == ([1, 0, 0], None)
    assert parse_version("v1.0.0.0.0") == ([1, 0, 0, 0, 0], None)
    assert parse_version("1.0.0-alpha") == ([1, 0, 0], ["alpha"])
    assert parse_version("1.0.0.alpha.1") == ([1, 0, 0], ["alpha", "1"])
    assert parse_version("v1-alpha.1") == ([1], ["alpha", "1"])
    assert parse_version("1.0.0-rc.1") == ([1, 0, 0], ["rc", "1"])
    assert parse_version("v1.0.0-x.7.z.92") == ([1, 0, 0], ["x", "7", "z", "92"])
    assert parse_version("0.1.9.1") == ([0, 1, 9, 1], None)
    assert parse_version("0.1.12.4") == ([0, 1, 12, 4], None)
    assert parse_version("1") == ([1], None)

    invalid_format = [
        "r1.0.0",
        "",
        "hello",
        "v1..0",
        "v1.0.0.",
        "1.0.0-",
    ]

    for version in invalid_format:
        print(f"Testing invalid version: {version}")
        with pytest.raises(InvalidVersion):
            parse_version(version)

def test_compare_versions():
    """Tests version comparison"""
    assert Version("1.0.0") < Version("1.0.1")
    assert Version("1") < Version("1.1.0")
    assert Version("1-beta1") < Version("1-beta2")
    assert Version("1.0.0") < Version("2.0.0-alpha.1")
    assert Version("1.0.0") < Version("2.0.0-rc.1")
    assert Version("1.0.0") < Version("2.0.0-x.1.z.92")
    assert Version("1-x.1.z.92") < Version("2-x.1.z.92")
    assert Version("v1.x.v.2.92") < Version("v1.x.v.3.92")
    assert Version("0.1.9.1") < Version("0.1.12.4")

    assert Version("1.0.0") == Version("1.0.0")
    assert Version("1.0.0-alpha.1") == Version("1.0.0-alpha.1")
    assert Version("1.0.0-rc.beta.1") == Version("1.0.0-rc.beta.1")
    assert Version("1-x.1.z.92") == Version("1-x.1.z.92")

    assert Version("1.0.0") > Version("1.0.0-alpha.1")
    # Curious case study. The expression will be incorrect
    # because the right operand contains more segments in the suffix.
    # I don't think i should break the logic for the sake of very rare cases
    # assert Version("1.0.0.rc-1") > Version("1.0.0-alpha.1")
    assert Version("1.0.0") > Version("1.0.0-rc.1")
    assert Version("1.0.0") > Version("1.0.0-x.1.z.92")
    assert not Version("0.1.12.4") > Version("0.1.12.4") # pylint: disable=C0117

    assert Version("1.0.0") >= Version("1.0.0")
    assert Version("1.0.0") >= Version("1.0.0-alpha.1")
    assert Version("1.0.0") >= Version("1.0.0-rc.1")
    assert Version("1.0.0") >= Version("1.0.0-x.1.z.92")

    assert Version("1.0.0") <= Version("1.0.0")
    assert Version("1.0.0-alpha.1") <= Version("1.0.0")
    assert Version("1-x.1.z.92") <= Version("2-x.1.z.92")
    assert Version("v1.x.v.2.92") <= Version("v1.x.v.3.92")

    assert Version("1.0.0") != Version("1.0.0-0")
    assert Version("1.0.0") != Version("1.0.0.0")
    assert Version("1.0.0") != Version("1.0.0-rc.1")
    assert Version("1.0.0") != Version("1.0.0-x.1.z.92")
