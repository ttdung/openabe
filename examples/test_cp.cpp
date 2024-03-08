/// 
/// Copyright (c) 2018 Zeutro, LLC. All rights reserved.
/// 
/// This file is part of Zeutro's OpenABE.
/// 
/// OpenABE is free software: you can redistribute it and/or modify
/// it under the terms of the GNU Affero General Public License as published by
/// the Free Software Foundation, either version 3 of the License, or
/// (at your option) any later version.
/// 
/// OpenABE is distributed in the hope that it will be useful,
/// but WITHOUT ANY WARRANTY; without even the implied warranty of
/// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
/// GNU Affero General Public License for more details.
/// 
/// You should have received a copy of the GNU Affero General Public
/// License along with OpenABE. If not, see <http://www.gnu.org/licenses/>.
/// 
/// You can be released from the requirements of the GNU Affero General
/// Public License and obtain additional features by purchasing a
/// commercial license. Buying such a license is mandatory if you
/// engage in commercial activities involving OpenABE that do not
/// comply with the open source requirements of the GNU Affero General
/// Public License. For more information on commerical licenses,
/// visit <http://www.zeutro.com>.
///
/// \brief  Example use of the OpenABE API with CP-ABE
///

//#include <iostream>
#include <fstream>
#include <string>
#include <cassert>
#include <openabe/openabe.h>
#include <openabe/zsymcrypto.h>

using namespace std;
using namespace oabe;
using namespace oabe::crypto;

int main(int argc, char **argv) {

  InitializeOpenABE();

  cout << "Testing CP-ABE context" << endl;

  OpenABECryptoContext cpabe("CP-ABE");

  string ct, pt1 = "hello world!", pt2, pt3;

  cpabe.generateParams();

  cpabe.keygen("|attr1|attr2", "key0");

  cpabe.encrypt("attr1 and attr2", pt1, ct);

  bool result = cpabe.decrypt("key0", ct, pt2);

  assert(result && pt1 == pt2);

  cout << "Recovered message: " << pt2 << endl;
  
  std::ifstream inFile;
    inFile.open("sample.json"); //open the input file

    std::stringstream strStream;
    strStream << inFile.rdbuf(); //read the file
    std::string str = strStream.str(); //str holds the content of the file

    inFile.close();

  cout << "sample file length: " << str.length()  << endl;
  // cout << "sample file: " << str << endl;

  // cout << "len of ct: " << ct.length() << endl;

  cpabe.encrypt("(attr1 and attr2) or attr3 or attr4 or attr5 or attr6 or attr7 or attr13 or attr14 or attr15 or attr16 or attr17 or attr61 or attr71 or attr131 or attr141 or attr151 or attr161 or attr171",str, ct);

  cout << "len of ct of sample file: " << ct.length() << endl;



  cpabe.keygen("|attr1|attr2", "key1");
  
  result = cpabe.decrypt("key1", ct, pt3);

  assert(result && str == pt3);

  cout << "Recovered message pt3 success: " << endl;

  ShutdownOpenABE();

  return 0;
}
