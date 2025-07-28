import 'package:flutter_riverpod/flutter_riverpod.dart';

T createAndRegister<T>(
  WidgetRef ref,
  T Function() creator,
  void Function(T) register,
) {
  final item = creator();
  register(item);
  return item;
}