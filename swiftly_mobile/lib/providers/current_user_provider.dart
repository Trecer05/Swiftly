import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../domain/user/models/user.dart';

final currentUserProvider = StateProvider<User?>((ref) => null);

